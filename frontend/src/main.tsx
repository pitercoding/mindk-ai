import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import App from "./App";

import { SelectedNoteProvider } from "@/context/SelectedNoteContext";

import "./index.css";

createRoot(
  document.getElementById("root")!,
).render(
  <StrictMode>

    <SelectedNoteProvider>

      <App />

    </SelectedNoteProvider>

  </StrictMode>,
);
