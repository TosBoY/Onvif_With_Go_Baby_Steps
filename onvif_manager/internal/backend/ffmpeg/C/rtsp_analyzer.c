/**
 * RTSP Stream Analyzer
 * 
 * This program connects to an RTSP stream and analyzes its resolution,
 * frame rate, bitrate, and codec information using FFmpeg libraries.
 */

#include <stdio.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

void display_stream_info(AVFormatContext *format_ctx, int stream_index) {
    AVStream *stream = format_ctx->streams[stream_index];
    AVCodecParameters *codec_params = stream->codecpar;
    
    // Skip non-video streams
    if (codec_params->codec_type != AVMEDIA_TYPE_VIDEO) {
        return;
    }
    
    double fps = 0;
    if (stream->avg_frame_rate.den && stream->avg_frame_rate.num) {
        fps = av_q2d(stream->avg_frame_rate);
    }
    
    const char *codec_name = avcodec_get_name(codec_params->codec_id);
    
    // Get bitrate (convert from bits/sec to kbps)
    int bitrate_kbps = 0;
    if (codec_params->bit_rate > 0) {
        bitrate_kbps = codec_params->bit_rate / 1000;
    }
      // Display all video stream parameters including bitrate
    printf("  Codec: %s\n", codec_name);
    printf("  Resolution: %dx%d\n", codec_params->width, codec_params->height);
    printf("  Frame rate: %.2f fps\n", fps);
    if (bitrate_kbps > 0) {
        printf("  Bitrate: %d kbps\n", bitrate_kbps);
    } else {
        printf("  Bitrate: Unknown\n");
    }
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        fprintf(stderr, "Usage: %s <rtsp_url>\n", argv[0]);
        return 1;
    }
    
    const char *rtsp_url = argv[1];
    AVFormatContext *format_ctx = NULL;
    int ret;
    
    // Register all codecs and formats (in newer FFmpeg versions this is unnecessary, but still safe)
    #if LIBAVCODEC_VERSION_INT < AV_VERSION_INT(58, 9, 100)
        av_register_all();
    #endif
    
    // Set up network
    avformat_network_init();
    
    // RTSP options - lower latency
    AVDictionary *options = NULL;
    // Tune for low latency
    av_dict_set(&options, "rtsp_transport", "tcp", 0); // Force TCP
    av_dict_set(&options, "max_delay", "500000", 0);   // 0.5s max delay
    av_dict_set(&options, "stimeout", "5000000", 0);   // 5s connection timeout
    
    printf("Connecting to: %s\n", rtsp_url);
    
    // Open input
    ret = avformat_open_input(&format_ctx, rtsp_url, NULL, &options);
    if (ret < 0) {
        char err_buf[AV_ERROR_MAX_STRING_SIZE];
        av_strerror(ret, err_buf, sizeof(err_buf));
        fprintf(stderr, "Could not open input: %s\n", err_buf);
        return 1;
    }
    
    // Retrieve stream information
    ret = avformat_find_stream_info(format_ctx, NULL);
    if (ret < 0) {
        char err_buf[AV_ERROR_MAX_STRING_SIZE];
        av_strerror(ret, err_buf, sizeof(err_buf));
        fprintf(stderr, "Could not find stream info: %s\n", err_buf);
        avformat_close_input(&format_ctx);
        return 1;
    }
      printf("\n===== RTSP Stream Analysis =====\n\n");
    
    // Find and display information for video and audio streams
    for (unsigned int i = 0; i < format_ctx->nb_streams; i++) {
        display_stream_info(format_ctx, i);
    }
    
    // Clean up
    avformat_close_input(&format_ctx);
    avformat_network_deinit();
    
    return 0;
}
