# Parrot versioning dashboard

Visualized dashboard that was developed in order to have overview of our translated projects in one table. We can choose and open any project and see available translation versions, add new one or delete some.

## Available Scripts

In the project directory, you can run:

### `pnpm start`

Runs the app in development mode.
Open [http://localhost:5174](http://localhost:5174) to view it in the browser.

The page will reload if you make edits, and you will see any lint errors in the console.

### `pnpm start:mocked`

Runs the app like `pnpm start`, but it sets the `VITE_MOCKED` variable to true to run mocked Axios requests.

### `pnpm lint`

Runs ESLint on the project with the rules specified in .eslintrc.json. The linter should be satisfied at all times!

### `pnpm build`

Builds the app for production to the `build` folder.
It optimizes the build for the best performance. The build is minified, and the filenames include hashes. Your app is ready to be deployed.

For more information on building with Vite, refer to the [Vite documentation](https://vitejs.dev/guide/build.html).
