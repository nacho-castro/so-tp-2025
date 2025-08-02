RUN = go run

# === Prueba: Corto Plazo ===
corto:
	nano ../kernel/configs/cortoplazo.json
	@echo "Comienza la prueba de Corto Plazo"
	$(RUN) kernel.go PLANI_CORTO_PLAZO 0 cortoplazo

# === Prueba: Mediano/Largo Plazo ===
lym:
	nano ../kernel/configs/medianoplazo.json
	@echo "Comienza la prueba de Mediano y Largo Plazo"
	$(RUN) kernel.go PLANI_LYM_PLAZO 0 medianoplazo

# === Prueba: SWAP ===
swap:
	nano ../kernel/configs/swap.json
	@echo "Comienza la prueba de SWAP"
	$(RUN) kernel.go MEMORIA_IO 90 swap

# === Prueba: CACHE ===
cache:
	nano ../kernel/configs/cache.json
	@echo "Comienza la prueba de Cache"
	$(RUN) kernel.go MEMORIA_BASE 256 cache

# === Prueba: TLB ===
tlb:
	nano ../kernel/configs/tlb.json
	@echo "Comienza la prueba de TLB"
	$(RUN) kernel.go MEMORIA_BASE_TLB 256 tlb


# === Prueba: Estabilidad General (EG) ===
eg:
	nano ../kernel/configs/estabilidad.json
	@echo "Comienza la prueba de Estabilidad"
	$(RUN) kernel.go ESTABILIDAD_GENERAL 0 estabilidad

# === Detener todos los procesos (por si quedan en background) ===
clean:
	nano ../kernel/configs/config.json
	@echo "Terminando módulo..."
	pkill -f kernel.go || true

