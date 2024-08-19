const esbuild = require("esbuild");
const { sassPlugin } = require("esbuild-sass-plugin");
const { sentryEsbuildPlugin } = require("@sentry/esbuild-plugin");

esbuild
    .context({
        sourcemap: true,
        entryPoints: ["frontend/index.tsx", "frontend/resources/styles/index.scss"],
        outdir: "public/assets",
        bundle: true,
        minify: true,

        plugins: [
            sassPlugin(),
            sentryEsbuildPlugin({
                authToken: process.env.SENTRY_AUTH_TOKEN,
                org: "fire-phoenix-games",
                project: "home-automation-frontend",
            }),
        ],
        loader: {
            ".png": "dataurl"
        }
    })
    .then((r) =>  {
            console.log("⚡ Build complete! ⚡");
            r.watch().then(r => console.log('watching...'));
    })
    .catch(() => process.exit(1));