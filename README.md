# go-jsonchanges

This packages finds changes between 2 json bodies and puts all changes in output json.
It's not a diffing tool, its more like collecting any change from A to B and even B to
A, favoring A to B changes first and putting them all together in output json. This may
feel a bit confusing cause there is nothing indicating the source of the change, but 
you get the changes in one look whether its coming from A or B. See the examples for
such beautiful confusions.
