NFC Key Reader for MicroNFCBoard
================================

This is micronfcboard firmware for handling NFC Lock keyfobs.

Planned Features
----------------

* Reading real UID from card
* USB NFC reader

SPI pins
--------

                SCLK            MOSI            MISO            SSEL
        SPI0    P0_6/P0_10      P0_9            P0_8            P0_2
                /P1_29
        SPI1    P1_15/P1_20     P0_21/P1_22     P0_22/P1_21     P1_19/P1_23
