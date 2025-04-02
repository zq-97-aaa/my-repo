package main

import (
	"fmt"
	"gorm.io/gorm"
	"helloworld/common"
	"log"
	"strings"
	"time"
)

var defaultUrl = []string{
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/abstract/abs_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i12.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Abstract/i9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f20.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f21.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f22.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f23.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f24.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f25.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f26.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f27.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f28.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f29.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f30.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f31.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f32.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f33.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f34.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f35.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f36.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f37.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f38.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Female/f9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m20.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m21.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m22.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m23.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m24.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m25.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m26.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m27.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m28.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m29.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m30.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m31.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m32.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m33.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m34.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m35.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m36.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m37.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m38.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m39.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m40.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m41.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Anime Male/m9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Female/f9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Cinematic Male/m9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/cute/cute_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i12.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i13.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i14.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i15.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i16.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i17.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i18.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i19.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i20.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i21.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i22.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Flat/i9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/futuristic/futu_9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/illustrative/illus_7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Impressionist/i9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i12.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i13.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i14.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i15.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i16.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i17.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i18.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i19.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i20.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i21.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i22.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i23.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i24.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Lofi/i9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i12.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i13.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i14.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i15.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i16.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i17.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i18.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i19.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Minimalist/i9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_12.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_13.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_14.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/minimalist/mini_9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n20.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n21.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n22.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n23.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n24.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n25.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n26.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n27.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n28.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n29.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n30.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n31.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n32.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n33.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n34.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n35.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n36.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n37.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n38.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n39.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n40.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n41.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n42.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n43.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n44.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n45.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n46.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Non Human/n9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f20.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f21.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f22.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f23.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f24.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f25.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f26.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f27.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f28.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f29.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f30.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f31.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f32.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f33.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f34.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f35.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f36.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f37.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Female/f9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m10.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m11.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m12.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m13.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m14.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m15.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m16.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m17.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m18.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m19.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m1.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m20.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m21.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m22.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m23.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m24.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m25.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m26.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m27.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m28.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m29.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m2.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m30.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m31.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m32.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m33.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m34.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m35.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m36.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m37.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m38.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m39.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m3.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m40.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m4.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m5.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m6.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m7.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m8.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/Realistic Male/m9.jpg",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/realistic/real_9.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_10.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_11.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_1.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_2.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_3.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_4.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_5.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_6.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_7.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_8.png",
	"https://flow-prompt-covers.s3.us-west-1.amazonaws.com/icon/vintage/vint_9.png",
}

func FilterPromptList2(user *TUserFeature, pRecall *PRecall) ([]*Prompt, error) {
	start := time.Now()
	defer func() {
		log.Println("FilterPromptList complete, costTime", time.Since(start).String())
	}()
	//db = InitDB()
	query := BuildPromptQuery(db, user)

	var prompts []*Prompt
	if err := query.Find(&prompts).Error; err != nil {
		return nil, fmt.Errorf("failed to query prompts: %w", err)
	}

	// 应用内存中的额外过滤（如曝光计数）
	filtered := make([]*Prompt, 0, len(prompts))
	for _, p := range prompts {
		if !RecallExposureLimit(user, p) {
			continue
		}
		filtered = append(filtered, p)
	}

	log.Printf("filtered %d prompts", len(filtered))
	return filtered, nil
}

func BuildPromptQuery(db *gorm.DB, user *TUserFeature) *gorm.DB {

	query := db.Model(&Prompt{}).Debug().Where("live = ?", true)

	// 1. PoolType过滤
	switch user.PoolType {
	case "sfw":
		query = query.Where("nsfw = ?", false)
	case "nsfw":
		query = query.Where("nsfw = ?", true)
	}

	// 2. 曝光限制过滤
	if len(user.LastRecommendedPrompts) > 0 {
		// 转换为SQL NOT IN条件
		excludedIDs := make([]string, 0, len(user.LastRecommendedPrompts))
		for id := range user.LastRecommendedPrompts {
			excludedIDs = append(excludedIDs, id)
		}
		query = query.Where("id NOT IN ?", excludedIDs)
	}

	// 3. 语言过滤
	if user.Language != "" {
		query = query.Where("language = ?", user.Language)
	}

	// 4. 标签匹配
	//if len(pRecall.Tags) > 0 {
	//	// 假设Tags是JSON数组，使用JSON_CONTAINS或类似函数
	//	// SQLite需要使用json_each处理
	//	for _, tagGroup := range pRecall.Tags {
	//		if len(tagGroup) == 0 {
	//			continue
	//		}
	//		query = query.Where(fmt.Sprintf(
	//			`EXISTS (SELECT 1 FROM json_each(tags) WHERE value IN (%s))`,
	//			buildPlaceholders(len(tagGroup)),
	//		), convertToInterfaceSlice(tagGroup)...)
	//	}
	//}

	// 5. 敏感图片过滤
	if user.DeviceType == common.MobileDeviceType && user.ReqScene == common.AppExplore {
		switch user.AB[common.AppSensitiveImageExp] {
		case common.ExpGroup1:
			//query = query.Where(`(
			//    json_extract(tagsMap, '$.sensitiveImage') IS NULL OR
			//    json_extract(tagsMap, '$.sensitiveImage') = 0
			//) AND (
			//    json_extract(tagsMap, '$.nvJianBeiYaoFu') IS NULL OR
			//    json_extract(tagsMap, '$.nvJianBeiYaoFu') = 0
			//) AND (
			//    json_extract(tagsMap, '$.seQing') IS NULL OR
			//    json_extract(tagsMap, '$.seQing') = 0
			//)`)
		case common.ExpGroup2:
			//query = query.Where(`(
			//    json_extract(tagsMap, '$.sensitiveImage2') IS NULL OR
			//    json_extract(tagsMap, '$.sensitiveImage2') = 0
			//)`)
		}
	}

	// 6. 时间过滤
	//if pRecall.CutoffDate != "" {
	//	cutoffTime, err := time.Parse(time.RFC3339, pRecall.CutoffDate)
	//	if err == nil {
	//		query = query.Where("created_at >= ?", cutoffTime)
	//	}
	//}

	// 7. 附加过滤
	//for filterK, filterV := range pRecall.AdditionalFilters {
	//	if len(filterV) == 0 {
	//		continue
	//	}
	//	switch filterK {
	//	case CategoryRecallFilter:
	//		if categoryId, err := strconv.ParseInt(filterV[0], 10, 64); err == nil {
	//			query = query.Where("category_id = ?", categoryId)
	//		}
	//	case AuthorRecallFilter:
	//		query = query.Where("user_id = ?", filterV[0])
	//	}
	//}

	// 8. 移动设备缩略图过滤
	if user.DeviceType == common.MobileDeviceType {
		query = query.Where("thumbnailUrl NOT IN ? ", defaultUrl)
	}

	return query
}

// 辅助函数：构建占位符字符串
func buildPlaceholders(n int) string {
	placeholders := make([]string, n)
	for i := 0; i < n; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

// 辅助函数：转换字符串切片为interface切片
func convertToInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
