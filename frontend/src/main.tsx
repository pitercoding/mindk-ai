import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import App from "./App";

import { SelectedKnowledgeProvider } from "@/context/SelectedKnowledgeContext";
import { ChatSessionProvider } from "@/context/ChatSessionContext";

import "./index.css";

createRoot(
  document.getElementById("root")!,
).render(
  <StrictMode>

    <SelectedKnowledgeProvider>

      <ChatSessionProvider>

        <App />

      </ChatSessionProvider>

    </SelectedKnowledgeProvider>

  </StrictMode>,
);
