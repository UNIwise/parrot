import React from "react";
import ReactDOM from "react-dom/client";
import { ReactQueryClientProvider } from "./api/client";
import { Parrot } from "./App";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ReactQueryClientProvider>
      <Parrot />
    </ReactQueryClientProvider>
  </React.StrictMode>,
);
