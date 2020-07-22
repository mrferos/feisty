package revisions

import (
	"github.com/r3labs/diff"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func RevisionWatchFilter(e event.UpdateEvent) bool {
	oldAnnotations := e.MetaOld.GetAnnotations()
	newAnnotations := e.MetaNew.GetAnnotations()

	if oldAnnotations == nil || newAnnotations == nil {
		return true
	}

	annotationDiff, err := diff.Diff(oldAnnotations, newAnnotations)
	if err != nil {
		return true
	}

	if len(annotationDiff) > 1 {
		return true
	}

	oldRevNumber := oldAnnotations[RevisionNumberAnnotation]
	newRevNumber := newAnnotations[RevisionNumberAnnotation]

	if oldRevNumber != newRevNumber {
		return false
	}

	return true
}
