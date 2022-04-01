import sveltePlugin from 'esbuild-svelte';
import sveltePreprocess from 'svelte-preprocess';
import { build } from 'esbuild';
import { existsSync, mkdirSync, rmSync } from 'fs';

// Parse arguments
const args = process.argv.slice(2);
const devMode = args.length > 0 && args[0] === 'dev';

// Make sure target directory exists
const outDir = './web/public/build';
if (!existsSync(outDir)) {
	mkdirSync(outDir);
}

// Clean up target directory
rmSync(outDir, { recursive: true });

// Prepare arguments
let watcher = false;
if (devMode) {
	console.log(`development mode started, watcher activated`);
	watcher = {
		onRebuild(error, result) {
			if (error) console.error('watch build failed:', error);
			else console.log('watch build succeeded:', result);
		},
	};
}

// Build the application
build({
	entryPoints: ['./web/src/app.ts'],
	outdir: outDir,
	minify: !devMode,
	bundle: true,
	sourcemap: devMode,
	watch: watcher,
	treeShaking: true,
	plugins: [
		sveltePlugin({
			cache: true,
			preprocess: sveltePreprocess(),
			compilerOptions: { dev: devMode },
		}),
	],
}).catch((err) => {
	console.error(err);
	process.exit(1);
});
