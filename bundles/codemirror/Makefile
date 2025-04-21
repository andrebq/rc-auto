.PHONY: bundle

bundle:
	npm i \
		codemirror \
		@codemirror/lang-javascript \
		@codemirror/lang-python \
		@codemirror/lang-go \
		@codemirror/legacy-modes \
		@codemirror/legacy-modes \
		@codemirror/theme-one-dark
	npm i \
		rollup \
		@rollup/plugin-node-resolve
	mkdir -p dist
	node_modules/.bin/rollup -c ./rollup.config.mjs