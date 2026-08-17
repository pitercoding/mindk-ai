import { SignIn } from "@clerk/react";

export default function AuthPage() {
    return (
        <main className="auth-page">

            <div className="auth-card">

                <div className="auth-brand">
                    <span className="auth-brand-mark" />
                    <span className="auth-brand-name">MindK AI</span>
                </div>

                <SignIn
                    withSignUp
                    forceRedirectUrl="/"
                    appearance={{
                        variables: {
                            colorPrimary: "#36e5ce",
                            colorBackground: "#141f30",
                            colorForeground: "#f4f8fb",
                            colorMutedForeground: "#a9b7c6",
                            colorInput: "rgba(15, 25, 40, 0.92)",
                            colorInputForeground: "#f4f8fb",
                            colorDanger: "#ff6172",
                            borderRadius: "8px",
                        },
                        elements: {
                            rootBox: "auth-clerk-root",
                            cardBox: "auth-clerk-card",
                        },
                    }}
                />

            </div>

        </main>
    );
}
