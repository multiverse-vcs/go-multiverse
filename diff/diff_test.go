package diff

import (
	"testing"
)

func TestMergeConflict(t *testing.T) {
	textO := `celery
garlic
onions
salmon
tomatoes
wine
`

	textA := `celery
salmon
tomatoes
garlic
onions
wine
`

	textB := `celery
garlic
salmon
tomatoes
onions
wine
`

	expect := `celery
salmon
tomatoes
garlic
<<<<<<<
onions
=======
salmon
tomatoes
onions
>>>>>>>
wine
`
	if Merge(textO, textA, textB) != expect {
		t.Errorf("unexpected merge result")
	}
}
