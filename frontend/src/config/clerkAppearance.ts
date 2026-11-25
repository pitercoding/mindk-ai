import type { ComponentProps } from "react";
import type { SignIn } from "@clerk/react";

type ClerkAppearance = NonNullable<ComponentProps<typeof SignIn>["appearance"]>;

/*
    Clerk theme aligned with the app palette (see :root in App.css).
    Shared by the sign-in page and the user menu so both look native.
*/
export const clerkAppearance: ClerkAppearance = {
    variables: {
        colorPrimary: "#36e5ce",
        colorBackground: "#141f30",
        colorForeground: "#f4f8fb",
        colorMutedForeground: "#a9b7c6",
        colorInput: "rgba(15, 25, 40, 0.92)",
        colorInputForeground: "#f4f8fb",
        colorDanger: "#ff6172",
        colorSuccess: "#3ddc97",
        colorWarning: "#f5b94a",
        borderRadius: "8px",
        fontFamily: "\"Inter Variable\", Inter, system-ui, sans-serif",
    },
    elements: {
        rootBox: "auth-clerk-root",
        cardBox: "auth-clerk-card",
        userButtonAvatarBox: "auth-user-avatar",
    },
};
