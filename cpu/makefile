RUN = go run

uno:
	nano ../cpu/configs/config1.json
	$(RUN) cpu.go 1

dos:
	nano ../cpu/configs/config2.json
	$(RUN) cpu.go 2

tres:
	nano ../cpu/configs/config3.json
	$(RUN) cpu.go 3

cuatro:
	nano ../cpu/configs/config4.json
	$(RUN) cpu.go 4

clean:
	@echo "Terminando módulos..."
	pkill -f cpu.go || true