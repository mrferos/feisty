package revisions

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/mrferos/feisty/api/v1alpha1"
	"github.com/mrferos/feisty/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

var RevisionNumberAnnotation = constants.FeistyAnnotationPrefix + "revision-number"

type Revision struct {
	client.Client
	Log logr.Logger
}

func (r *Revision) structMd5(obj interface{}) (string, error) {
	keyValueJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	md5bytes := md5.Sum([]byte(keyValueJson))

	return fmt.Sprintf("%x", md5bytes), nil
}

func (r *Revision) CreateIfNeeded(appName types.NamespacedName, ctx context.Context) error {
	log := r.Log.WithValues("source", "revision", "appName", appName.Name, "appNamespace", appName.Namespace)

	var app v1alpha1.Application
	var cfg v1alpha1.ApplicationConfig

	if err := r.Get(ctx, appName, &app); err != nil {
		log.Error(err, "Unable to fetch Application")
		return err
	}

	// The app may not have any config so we'll ignore a not found error
	if err := r.Get(ctx, appName, &cfg); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "Unable to fetch ApplicationConfig")
			return err
		} else {
			log.Info("Unable to fetch ApplicationConfig", "configName", appName)
		}
	}

	currentRevisionNumber := 0
	var err error
	if app.Annotations != nil {
		if val, ok := app.Annotations[RevisionNumberAnnotation]; ok {
			currentRevisionNumber, err = strconv.Atoi(val)
			if err != nil {
				log.Error(err, "There was an error parsing the revision number from the application")
				return err
			}
		}
	}

	cfgHash := ""
	appHash := ""
	revNumber := currentRevisionNumber + 1
	revName := app.Name + "-v" + strconv.Itoa(revNumber)
	prevRevName := ""

	if currentRevisionNumber > 0 {
		prevRevName = app.Name + "-v" + strconv.Itoa(currentRevisionNumber)
	}

	appHash, err = r.structMd5(app.Spec)
	if err != nil {
		log.Error(err, "Could not hash the Application")
		return err
	}

	if cfg.Name != "" {
		cfgHash, err = r.structMd5(cfg.Spec)
	}

	rev := v1alpha1.ApplicationRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:      revName,
			Namespace: app.Namespace,
		},
		Spec: v1alpha1.ApplicationRevisionSpec{
			App:     app.Spec,
			Cfg:     cfg.Spec,
			AppHash: appHash,
			CfgHash: cfgHash,
		},
	}

	prevRev := v1alpha1.ApplicationRevision{}
	prevNamespacedName := types.NamespacedName{
		Namespace: app.Namespace,
		Name:      prevRevName,
	}

	// We may have deleted the prevision religion, so we'll ignore a not found error
	if prevRevName != "" {
		if err := r.Get(ctx, prevNamespacedName, &prevRev); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "Unable to fetch previous ApplicationRevision")
				return err
			} else {
				log.Info("Unable to fetch previous ApplicationRevision", "prevName", prevRevName)
			}
		}
	}

	// If the hashes don't match then we'll save the new revision
	if rev.Spec.AppHash != prevRev.Spec.AppHash || rev.Spec.CfgHash != prevRev.Spec.CfgHash {
		log.Info("Saving new revision ", "revisionName", revName)
		if err := r.Create(ctx, &rev); err != nil {
			log.Error(err, "Could not create revision", "revisionName", revName)
			return err
		}

		if app.Annotations == nil {
			app.Annotations = map[string]string{}
		}

		app.Annotations[RevisionNumberAnnotation] = strconv.Itoa(revNumber)
		if err := r.Update(ctx, &app); err != nil {
			log.Error(err, "Could not update application with current revision number")
			return err
		}
	}

	return nil
}
