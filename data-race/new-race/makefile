COMPILER=g++
CFLAGS=-Wall
EXECUTABLE=sets_test

SOURCES=main.cpp lazySync.h set.h fineGrainedSync.h

all: $(SOURCES)
	$(COMPILER) $(CFLAGS) -o $(EXECUTABLE) $(SOURCES) -pthread -g
all-with-sanitizer:$(SOURCES)
	$(COMPILER) $(CFLAGS) -o $(EXECUTABLE) $(SOURCES) -pthread -fsanitize=thread -g
clean:
	rm -f *.o $(EXECUTABLE)
