include $(GOROOT)/src/Make.inc

TARG=pokerodds
GOFILES= \
	card.go \
	cardBag.go \
	hand.go \
	subsetChooser.go \

poker-odds: _obj/$(TARG).a
	$(GC) -I_obj poker-odds.go
	$(LD) -L_obj -o $@ poker-odds.$O
	@echo "Done. Executable is $@"

include $(GOROOT)/src/Make.pkg
