# Makefile for RTSP Analyzer

CC = gcc
CFLAGS = -Wall -Wextra -O2
LDFLAGS = -lavformat -lavcodec -lavutil -lswscale -lavdevice -lswresample

TARGET = rtsp_analyzer
SRC = rtsp_analyzer.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS)

clean:
	del /Q $(TARGET).exe

.PHONY: all clean
