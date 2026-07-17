import { NavLink } from "react-router-dom";

const menuItems = [
    {
        label: "Dashboard",
        path: "/",
        icon: "⌂",
    },
    {
        label: "All Notes",
        path: "/notes",
        icon: "▤",
    },
    {
        label: "Documents",
        path: "/documents",
        icon: "□",
    },
    {
        label: "Links",
        path: "/links",
        icon: "∞",
    },
    {
        label: "PDFs",
        path: "/pdfs",
        icon: "⌑",
        badge: "Soon",
    },
    {
        label: "Settings",
        path: "/settings",
        icon: "⚙",
    },
];

export default function Sidebar() {
    return (
        <aside className="sidebar">

            <header className="sidebar-header">
                <h1>
                    🧠 MindK AI
                </h1>
            </header>


            <nav className="sidebar-menu">

                {menuItems.map((item) => (
                    <NavLink
                        key={item.label}
                        data-icon={item.icon}
                        to={item.path}
                        className={({ isActive }) =>
                            isActive
                                ? "sidebar-link active"
                                : "sidebar-link"
                        }
                    >
                        <span>
                            {item.label}
                        </span>

                        {item.badge && (
                            <small>
                                {item.badge}
                            </small>
                        )}

                    </NavLink>
                ))}

            </nav>


            <footer className="sidebar-footer">

                <button>
                    ⬆ Upload
                </button>

            </footer>

        </aside>
    );
}
