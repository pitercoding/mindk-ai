import { UserButton } from "@clerk/react";

export default function Header() {
    return (
        <header className="app-header">

            <div className="search-container">

                <input
                    type="text"
                    placeholder="Search your notes, links..."
                />

            </div>


            <div className="header-actions">

                <button
                    className="notification-button"
                    type="button"
                    aria-label="Notifications"
                    title="Coming soon"
                    disabled
                >
                    🔔
                </button>

                <UserButton
                    appearance={{
                        elements: {
                            userButtonAvatarBox: "auth-user-avatar",
                        },
                    }}
                />

            </div>

        </header>
    );
}
