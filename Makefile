BUILD_DIR_WINDOWS=.\cmd\go-chat\build
WORK_DIR_WINDOWS=.\cmd\go-chat
CONFIG_DIR_WINDOWS=.\cmd\go-chat\config

run.windows:
	go run $(WORK_DIR_WINDOWS)\. \
		-app.config.files $(CONFIG_DIR_WINDOWS)\application.yaml \
		-logger.config.file $(CONFIG_DIR_WINDOWS)\logger.json \
		-env.vars.file $(CONFIG_DIR_WINDOWS)\sample.env
