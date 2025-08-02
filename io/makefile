RUN = go run

uno:
	$(RUN) io.go

io:
	nano ../io/configs/config.json
	$(RUN) io.go

# === Prueba: Corto Plazo ===
corto:
	$(RUN) io.go DISCO1

# === Prueba: Mediano/Largo Plazo ===
lym:
	$(RUN) io.go DISCO1

# === Prueba: SWAP ===
swap:
	$(RUN) io.go DISCO1

# === Prueba: CACHE ===
cache:
	$(RUN) io.go DISCO1

# === Prueba: Estabilidad General (EG) ===
eg:
	$(RUN) io.go DISCO1 &
	$(RUN) io.go DISCO2 &
	$(RUN) io.go DISCO3 &
	$(RUN) io.go DISCO4

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulo..."
	pkill -f io.go || true
