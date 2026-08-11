import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import App from "./App";

import { SelectedNoteProvider } from "@/context/SelectedNoteContext";
import { ChatSessionProvider } from "@/context/ChatSessionContext";

import "./index.css";

createRoot(
  document.getElementById("root")!,
).render(
  <StrictMode>

    <SelectedNoteProvider>

      <ChatSessionProvider>

        <App />

      </ChatSessionProvider>

    </SelectedNoteProvider>

  </StrictMode>,
);
