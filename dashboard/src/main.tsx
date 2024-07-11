import React from "react";
import ReactDOM from "react-dom/client";

import { Box, CssBaseline, CssVarsProvider } from "@mui/joy";
import { BrowserRouter } from "react-router-dom";
import "./index.css";
import Routes from "./routes";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <CssVarsProvider disableTransitionOnChange>
      <CssBaseline />
      <Box sx={{ display: "flex", minHeight: "100dvh" }}>
        <Box
          component="main"
          sx={{
            px: { xs: 2, md: 10 },
            pt: {
              xs: "calc(12px + var(--Header-height))",
              sm: "calc(12px + var(--Header-height))",
              md: 3,
            },
            pb: { xs: 2, sm: 2, md: 3 },
            flex: 1,
            display: "flex",
            flexDirection: "column",
            minWidth: 0,
            height: "100dvh",
            gap: 1,
          }}
        >
          <BrowserRouter>
            <Routes />
          </BrowserRouter>
        </Box>
      </Box>
    </CssVarsProvider>
  </React.StrictMode>,
);
