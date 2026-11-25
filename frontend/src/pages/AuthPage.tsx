import { SignIn } from "@clerk/react";

import BrandMark from "@/components/common/BrandMark";

import { clerkAppearance } from "@/config/clerkAppearance";

export default function AuthPage() {
    return (
        <main className="auth-page">

            <div className="auth-card">

                <div className="auth-brand">
                    <BrandMark size={32} />
                    <span className="auth-brand-name">MindK AI</span>
                </div>

                <p className="auth-tagline">
                    Your second brain: notes and documents, searchable with AI.
                </p>

                <SignIn
                    withSignUp
                    forceRedirectUrl="/"
                    appearance={clerkAppearance}
                />

            </div>

        </main>
    );
}
