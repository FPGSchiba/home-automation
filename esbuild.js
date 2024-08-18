const esbuild = require("esbuild");
const { sassPlugin } = require("esbuild-sass-plugin");

esbuild
    .build({
        entryPoints: ["frontend/index.tsx", "frontend/resources/styles/index.scss"],
        outdir: "public/assets",
        bundle: true,
        watch: true,
        minify: true,
        plugins: [sassPlugin()],
    })
    .then(() => console.log("⚡ Build complete! ⚡"))
    .catch(() => process.exit(1));