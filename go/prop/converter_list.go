package prop

// converters contain all the property converters.
var converters = []*SGFConverter{
	sizeConv,
	placementsConv,
	movesConv,
	komiConv,
	initPlayerConv,
	commentConv,
}

var propToConv = func(conv []*SGFConverter) map[Prop]*SGFConverter {
	mp := make(map[Prop]*SGFConverter)
	for _, c := range conv {
		for _, p := range c.Props {
			mp[p] = c
		}
	}
	return mp
}(converters)
