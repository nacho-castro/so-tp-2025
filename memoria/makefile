RUN = go run

memoria:
	$(RUN) memoria.go

# === Prueba: Corto Plazo ===
corto:
	nano ../memoria/configs/corto.json
	@echo "Comienza la prueba de Corto Plazo"
	$(RUN) memoria.go corto

# === Prueba: Mediano/Largo Plazo ===
lym:
	nano ../memoria/configs/medianoLargo.json
	@echo "Comienza la prueba de Mediano y Largo plazo"
	$(RUN) memoria.go medianoLargo

# === Prueba: SWAP ===
swapp:
	nano ../memoria/configs/swap.json
	@echo "Comienza la prueba de SWAP"
	$(RUN) memoria.go swap

# === Prueba: CACHE ===
cache:
	nano ../memoria/configs/cache.json
	@echo "Comienza la prueba de Cache"
	$(RUN) memoria.go cache

# === Prueba: TLB ===
tlb:
	nano ../memoria/configs/tlb.json
	@echo "Comienza la prueba de TLB"
	$(RUN) memoria.go tlb

# === Prueba: Estabilidad General (EG) ===
eg:
	nano ../memoria/configs/estabilidad.json
	@echo "Comienza la prueba de Estabilidad General"
	$(RUN) memoria.go estabilidad

# === Detener todos los procesos (por si quedan en background) ===
clean:
	@echo "Terminando módulo..."
	pkill -f memoria.go || true