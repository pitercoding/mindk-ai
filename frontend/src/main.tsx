import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ClerkProvider } from "@clerk/react";

import App from "./App";

import { SelectedKnowledgeProvider } from "@/context/SelectedKnowledgeContext";
import { ChatSessionProvider } from "@/context/ChatSessionContext";
import ToastProvider from "@/components/common/ToastProvider";

import "@fontsource-variable/inter";
import "./index.css";

const publishableKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

if (!publishableKey) {
  throw new Error("Missing VITE_CLERK_PUBLISHABLE_KEY");
}

createRoot(
  document.getElementById("root")!,
).render(
  <StrictMode>

    <ClerkProvider publishableKey={publishableKey} afterSignOutUrl="/auth">

      <ToastProvider>

        <SelectedKnowledgeProvider>

          <ChatSessionProvider>

            <App />

          </ChatSessionProvider>

        </SelectedKnowledgeProvider>

      </ToastProvider>

    </ClerkProvider>

  </StrictMode>,
);
