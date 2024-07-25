import React from "react";
import ReactDOM from "react-dom/client";

import { Container, CssBaseline, CssVarsProvider } from "@mui/joy";
import { BrowserRouter } from "react-router-dom";
import { ReactQueryClientProvider } from "./api/client";
import "./index.scss";
import Routes from "./routes";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <CssVarsProvider disableTransitionOnChange>
      <CssBaseline />
      <Container
        component="main"
        maxWidth="xl"
        sx={{
          display: "flex",
          flexDirection: "column",
          gap: 1,
        }}
      >
        <ReactQueryClientProvider>
          <BrowserRouter>
            <Routes />
          </BrowserRouter>
        </ReactQueryClientProvider>
      </Container>
    </CssVarsProvider>
  </React.StrictMode>,
);
