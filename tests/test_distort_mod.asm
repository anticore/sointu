%define BPM 100

%include "../src/sointu_header.inc"

BEGIN_PATTERNS
    PATTERN 64, HLD, HLD, HLD, HLD, HLD, HLD, HLD,  0, 0, 0, 0, 0, 0, 0, 0
END_PATTERNS

BEGIN_TRACKS
    TRACK VOICES(1),0
END_TRACKS

BEGIN_PATCH
    BEGIN_INSTRUMENT VOICES(1) ; Instrument0
        SU_ENVELOPE MONO,ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(80),GAIN(128)
        SU_DISTORT  MONO,DRIVE(32)
        SU_ENVELOPE MONO, ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(80),GAIN(128)
        SU_DISTORT  MONO, DRIVE(96)
        SU_OSCILLAT MONO,TRANSPOSE(70),DETUNE(64),PHASE(64),COLOR(128),SHAPE(64),GAIN(128),FLAGS(SINE+LFO)
        SU_SEND     MONO,AMOUNT(68),LOCALPORT(1,0)
        SU_SEND     MONO,AMOUNT(68),LOCALPORT(3,0) + SEND_POP
        SU_OUT      STEREO,GAIN(128)
    END_INSTRUMENT
END_PATCH

%include "../src/sointu_footer.inc"
