package arpcdata

type ProtoMapper[Proto comparable, Entity comparable] map[Proto]Entity

type ProtoConverter[Proto comparable, Entity comparable] interface {
	ToProto(src Entity) Proto
	FromProto(src Proto) Entity
}

type protoConverterImpl[Proto comparable, Entity comparable] struct {
	mapper        ProtoMapper[Proto, Entity]
	protoDefault  Proto
	entityDefault Entity
}

func (c *protoConverterImpl[Proto, Entity]) ToProto(src Entity) Proto {
	for proto, entity := range c.mapper {
		if entity == src {
			return proto
		}
	}

	return c.protoDefault
}

func (c *protoConverterImpl[Proto, Entity]) FromProto(src Proto) Entity {
	if entity, ok := c.mapper[src]; ok {
		return entity
	}

	return c.entityDefault
}

func NewProtoConverter[Proto comparable, Entity comparable](
	mapper ProtoMapper[Proto, Entity],
	protoDefault Proto,
	entityDefault Entity,
) ProtoConverter[Proto, Entity] {
	return &protoConverterImpl[Proto, Entity]{
		mapper:        mapper,
		protoDefault:  protoDefault,
		entityDefault: entityDefault,
	}
}
