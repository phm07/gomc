package registry

import (
	_ "embed"
	"encoding/json"
	"gomc/src/nbt"
)

//go:embed registry.json
var registryJson []byte

var Registry struct {
	MinecraftTrimPattern struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				TemplateItem string `json:"template_item"`
				Description  struct {
					Translate string `json:"translate"`
				} `json:"description"`
				AssetId string `json:"asset_id"`
				Decal   byte   `json:"decal"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:trim_pattern"`
	MinecraftTrimMaterial struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				Ingredient     string  `json:"ingredient"`
				AssetName      string  `json:"asset_name"`
				ItemModelIndex float32 `json:"item_model_index"`
				Description    struct {
					Color     string `json:"color"`
					Translate string `json:"translate"`
				} `json:"description"`
				OverrideArmorMaterials struct {
					Diamond   string `json:"diamond,omitempty"`
					Gold      string `json:"gold,omitempty"`
					Iron      string `json:"iron,omitempty"`
					Netherite string `json:"netherite,omitempty"`
				} `json:"override_armor_materials,omitempty"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:trim_material"`
	MinecraftChatType struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				Chat struct {
					TranslationKey string   `json:"translation_key"`
					Parameters     []string `json:"parameters"`
					Style          struct {
						Color  string `json:"color"`
						Italic byte   `json:"italic"`
					} `json:"style,omitempty"`
				} `json:"chat"`
				Narration struct {
					TranslationKey string   `json:"translation_key"`
					Parameters     []string `json:"parameters"`
				} `json:"narration"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:chat_type"`
	MinecraftDimensionType struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				PiglinSafe                  byte    `json:"piglin_safe"`
				Natural                     byte    `json:"natural"`
				AmbientLight                float32 `json:"ambient_light"`
				MonsterSpawnBlockLightLimit int     `json:"monster_spawn_block_light_limit"`
				Infiniburn                  string  `json:"infiniburn"`
				RespawnAnchorWorks          byte    `json:"respawn_anchor_works"`
				HasSkylight                 byte    `json:"has_skylight"`
				BedWorks                    byte    `json:"bed_works"`
				Effects                     string  `json:"effects"`
				HasRaids                    byte    `json:"has_raids"`
				LogicalHeight               int     `json:"logical_height"`
				CoordinateScale             float64 `json:"coordinate_scale"`
				MonsterSpawnLightLevel      any     `json:"monster_spawn_light_level"`
				MinY                        int     `json:"min_y"`
				Ultrawarm                   byte    `json:"ultrawarm"`
				HasCeiling                  byte    `json:"has_ceiling"`
				Height                      int     `json:"height"`
				FixedTime                   int64   `json:"fixed_time,omitempty"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:dimension_type"`
	MinecraftDamageType struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				Scaling          string  `json:"scaling"`
				Exhaustion       float32 `json:"exhaustion"`
				MessageId        string  `json:"message_id"`
				DeathMessageType string  `json:"death_message_type,omitempty"`
				Effects          string  `json:"effects,omitempty"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:damage_type"`
	MinecraftWorldgenBiome struct {
		Type  string `json:"type"`
		Value []struct {
			Name    string `json:"name"`
			Id      int    `json:"id"`
			Element struct {
				Effects struct {
					Music struct {
						ReplaceCurrentMusic byte   `json:"replace_current_music"`
						MaxDelay            int    `json:"max_delay"`
						Sound               string `json:"sound"`
						MinDelay            int    `json:"min_delay"`
					} `json:"music,omitempty"`
					SkyColor      int `json:"sky_color"`
					GrassColor    int `json:"grass_color,omitempty"`
					FoliageColor  int `json:"foliage_color,omitempty"`
					WaterFogColor int `json:"water_fog_color"`
					FogColor      int `json:"fog_color"`
					WaterColor    int `json:"water_color"`
					MoodSound     struct {
						TickDelay         int     `json:"tick_delay"`
						Offset            float64 `json:"offset"`
						Sound             string  `json:"sound"`
						BlockSearchExtent int     `json:"block_search_extent"`
					} `json:"mood_sound"`
					AmbientSound   string `json:"ambient_sound,omitempty"`
					AdditionsSound struct {
						Sound      string  `json:"sound"`
						TickChance float64 `json:"tick_chance"`
					} `json:"additions_sound,omitempty"`
					Particle struct {
						Probability float32 `json:"probability"`
						Options     struct {
							Type string `json:"type"`
						} `json:"options"`
					} `json:"particle,omitempty"`
					GrassColorModifier string `json:"grass_color_modifier,omitempty"`
				} `json:"effects"`
				HasPrecipitation    byte    `json:"has_precipitation"`
				Temperature         float32 `json:"temperature"`
				Downfall            float32 `json:"downfall"`
				TemperatureModifier string  `json:"temperature_modifier,omitempty"`
			} `json:"element"`
		} `json:"value"`
	} `json:"minecraft:worldgen/biome"`
}

var (
	RegistryNBT      nbt.Tag
	RegistryNBTBytes []byte
)

func init() {
	err := json.Unmarshal(registryJson, &Registry)
	if err != nil {
		panic(err)
	}

	for i, v := range Registry.MinecraftDimensionType.Value {
		switch l := v.Element.MonsterSpawnLightLevel.(type) {
		case float64:
			Registry.MinecraftDimensionType.Value[i].Element.MonsterSpawnLightLevel = int(l)
		case map[string]any:
			m := l["value"]
			if m, ok := m.(map[string]any); ok {
				for key, val := range m {
					if val, ok := val.(float64); ok {
						m[key] = int(val)
					}
				}
			}
		}
	}

	RegistryNBT = nbt.Marshal(Registry)
	RegistryNBTBytes = RegistryNBT.Marshal()
}
