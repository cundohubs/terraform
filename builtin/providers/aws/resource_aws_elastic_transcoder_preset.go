package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElasticTranscoderPreset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticTranscoderPresetCreate,
		Read:   resourceAwsElasticTranscoderPresetRead,
		Update: resourceAwsElasticTranscoderPresetUpdate,
		Delete: resourceAwsElasticTranscoderPresetDelete,

		Schema: map[string]*schema.Schema{
			"audio": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					// elastictranscoder.AudioParameters
					Schema: map[string]*schema.Schema{
						"audio_packing_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"bitrate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"channels": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec_options": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sample_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"container": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"thumbnails": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					// elastictranscoder.Thumbnails
					Schema: map[string]*schema.Schema{
						"aspect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_height": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_width": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"padding_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"resolution:": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sizing_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"video": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					// elastictranscoder.VideoParameters
					Schema: map[string]*schema.Schema{
						"aspect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"bitrate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec_options": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
						},
						"display_apect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fixed_gop": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"frame_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_frame_max_dist": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_frame_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_height": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_width": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"padding_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"resolution": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sizing_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"watermarks": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								// elastictranscoder.PresetWatermark
								Schema: map[string]*schema.Schema{
									"horizontal_align": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"horizaontal_offset": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"max_height": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"max_width": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"opacity": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"sizing_policy": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"target": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"vertical_align": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"vertical_offset": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsElasticTranscoderPresetCreate(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	req := &elastictranscoder.CreatePresetInput{
		Audio:       expandETAudioParams(d),
		Container:   aws.String(d.Get("container").(string)),
		Description: getStringPtr(d, "description"),
		Name:        aws.String(d.Get("name").(string)),
		Thumbnails:  expandETThumbnails(d),
		Video:       exapndETVideoParams(d),
	}

	log.Printf("[DEBUG] Elastic Transcoder Preset create opts: %s", req)
	resp, err := elastictranscoderconn.CreatePreset(req)
	if err != nil {
		return fmt.Errorf("Error creating Elastic Transcoder Preset: %s", err)
	}

	if resp.Warning != nil {
		log.Printf("[WARN] Elastic Transcoder Preset: %s", *resp.Warning)
	}

	return resourceAwsElasticTranscoderPipelineUpdate(d, meta)

}

func expandETThumbnails(d *schema.ResourceData) *elastictranscoder.Thumbnails {
	set, ok := d.GetOk("thumbnails")
	if !ok {
		return nil
	}

	s := set.(*schema.Set)
	if s == nil || s.Len() == 0 {
		return nil
	}
	t := s.List()[0].(map[string]interface{})

	return &elastictranscoder.Thumbnails{
		AspectRatio:   getStringPtr(t, "aspect_ratio"),
		Format:        getStringPtr(t, "format"),
		Interval:      getStringPtr(t, "interval"),
		MaxHeight:     getStringPtr(t, "max_height"),
		MaxWidth:      getStringPtr(t, "max_width"),
		PaddingPolicy: getStringPtr(t, "padding_policy"),
		Resolution:    getStringPtr(t, "resolution"),
		SizingPolicy:  getStringPtr(t, "sizing_policy"),
	}
}

func expandETAudioParams(d *schema.ResourceData) *elastictranscoder.AudioParameters {
	set, ok := d.GetOk("audio")
	if !ok {
		return nil
	}

	s := set.(*schema.Set)
	if s == nil || s.Len() == 0 {
		return nil
	}
	audio := s.List()[0].(map[string]interface{})

	return &elastictranscoder.AudioParameters{
		AudioPackingMode: getStringPtr(audio, "audio_packing_mode"),
		BitRate:          getStringPtr(audio, "bitrate"),
		Channels:         getStringPtr(audio, "channels"),
		Codec:            getStringPtr(audio, "codec"),
		CodecOptions:     expandETAudioCodecOptions(audio["codec_options"].(*schema.Set)),
		SampleRate:       getStringPtr(audio, "sample_rate"),
	}
}

func expandETAudioCodecOptions(s *schema.Set) *elastictranscoder.AudioCodecOptions {
	if s == nil || s.Len() == 0 {
		return nil
	}

	codec := s.List()[0].(map[string]interface{})

	codecOpts := &elastictranscoder.AudioCodecOptions{
		BitDepth: getStringPtr(codec, "bit_depth"),
		BitOrder: getStringPtr(codec, "bit_prder"),
		Profile:  getStringPtr(codec, "profile"),
		Signed:   getStringPtr(codec, "signed"),
	}

	return codecOpts
}

func exapndETVideoParams(d *schema.ResourceData) *elastictranscoder.VideoParameters {
	set, ok := d.GetOk("video")
	if !ok {
		return nil
	}

	s := set.(*schema.Set)
	if s == nil || s.Len() == 0 {
		return nil
	}
	p := s.List()[0].(map[string]interface{})

	return &elastictranscoder.VideoParameters{
		AspectRatio:        getStringPtr(p, "aspect_ratio"),
		BitRate:            getStringPtr(p, "bit_rate"),
		Codec:              getStringPtr(p, "codec"),
		CodecOptions:       stringMapToPointers(p["codec_options"].(map[string]interface{})),
		DisplayAspectRatio: getStringPtr(p, "display_aspect_ratio"),
		FixedGOP:           getStringPtr(p, "fixed_gop"),
		FrameRate:          getStringPtr(p, "frame_rate"),
		KeyframesMaxDist:   getStringPtr(p, "key_frame_max_dist"),
		MaxFrameRate:       getStringPtr(p, "max_frame_rate"),
		MaxHeight:          getStringPtr(p, "max_height"),
		MaxWidth:           getStringPtr(p, "max_width"),
		PaddingPolicy:      getStringPtr(p, "padding_policy"),
		Resolution:         getStringPtr(p, "resolution"),
		SizingPolicy:       getStringPtr(p, "sizing_policy"),
		Watermarks:         expandETWatermarks(p["watermarks"].(*schema.Set)),
	}
}

func expandETWatermarks(s *schema.Set) []*elastictranscoder.PresetWatermark {
	var watermarks []*elastictranscoder.PresetWatermark

	for _, w := range s.List() {
		watermark := &elastictranscoder.PresetWatermark{
			HorizontalAlign:  getStringPtr(w, "horizontal_align"),
			HorizontalOffset: getStringPtr(w, "horizontal_offset"),
			Id:               getStringPtr(w, "id"),
			MaxHeight:        getStringPtr(w, "max_height"),
			MaxWidth:         getStringPtr(w, "max_width"),
			Opacity:          getStringPtr(w, "opacity"),
			SizingPolicy:     getStringPtr(w, "sizing_policy"),
			Target:           getStringPtr(w, "target"),
			VerticalAlign:    getStringPtr(w, "vertical_align"),
			VerticalOffset:   getStringPtr(w, "vertical_offset"),
		}
		watermarks = append(watermarks, watermark)
	}

	return watermarks
}

func resourceAwsElasticTranscoderPresetUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticTranscoderPresetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticTranscoderPresetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
