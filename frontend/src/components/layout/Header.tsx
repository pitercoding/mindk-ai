import { UserButton } from "@clerk/react";

interface HeaderProps {
    isMenuOpen: boolean;
    onMenuClick: () => void;
}

export default function Header({
    isMenuOpen,
    onMenuClick,
}: HeaderProps) {
    return (
        <header className="app-header">

            <button
                type="button"
                className="menu-button"
                aria-label={isMenuOpen ? "Close menu" : "Open menu"}
                aria-expanded={isMenuOpen}
                aria-controls="app-sidebar"
                onClick={onMenuClick}
            >
                <span aria-hidden="true" />
            </button>

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
