import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { NavLink, useNavigate } from "react-router-dom";

import { useChatSession } from "@/context/ChatSessionContext";

const menuItems = [
    {
        label: "Dashboard",
        path: "/",
        icon: "⌂",
        end: true,
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

interface MenuState {
    sessionId: number;
    top: number;
    left: number;
}

export default function Sidebar() {

    const navigate = useNavigate();

    const {
        sessions,
        currentSessionId,
        newChat,
        selectSession,
        renameSession,
        deleteSession,
    } = useChatSession();

    const [menu, setMenu] = useState<MenuState | null>(null);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [editingTitle, setEditingTitle] = useState("");

    useEffect(() => {

        if (!menu) {
            return;
        }

        function closeMenu() {
            setMenu(null);
        }

        document.addEventListener("mousedown", closeMenu);
        window.addEventListener("resize", closeMenu);
        window.addEventListener("scroll", closeMenu, true);

        return () => {
            document.removeEventListener("mousedown", closeMenu);
            window.removeEventListener("resize", closeMenu);
            window.removeEventListener("scroll", closeMenu, true);
        };

    }, [menu]);

    function handleNewChat() {
        setMenu(null);
        newChat();
        navigate("/");
    }

    async function handleSelect(id: number) {
        setMenu(null);
        await selectSession(id);
        navigate("/");
    }

    function toggleMenu(
        event: React.MouseEvent<HTMLButtonElement>,
        sessionId: number,
    ) {

        if (menu?.sessionId === sessionId) {
            setMenu(null);
            return;
        }

        const rect = event.currentTarget.getBoundingClientRect();
        const menuWidth = 120;

        setMenu({
            sessionId,
            top: rect.bottom + 4,
            left: Math.min(
                rect.right - menuWidth,
                window.innerWidth - menuWidth - 8,
            ),
        });
    }

    function startRename(id: number, title: string) {
        setMenu(null);
        setEditingId(id);
        setEditingTitle(title);
    }

    async function commitRename(id: number) {

        const title = editingTitle.trim();

        setEditingId(null);

        if (!title) {
            return;
        }

        try {
            await renameSession(id, title);
        } catch (error) {
            console.error("Failed to rename conversation:", error);
        }
    }

    async function handleDelete(id: number) {

        setMenu(null);

        if (!window.confirm("Delete this conversation?")) {
            return;
        }

        try {
            await deleteSession(id);
        } catch (error) {
            console.error("Failed to delete conversation:", error);
        }
    }

    const activeMenuSession =
        menu && sessions.find(session => session.id === menu.sessionId);

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
                        end={item.end}
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


            <div className="sidebar-conversations">

                <div className="sidebar-section-title">
                    Conversations
                </div>

                <button
                    className="sidebar-new-chat-button"
                    onClick={handleNewChat}
                >
                    + New Chat
                </button>

                <ul className="conversation-list">

                    {sessions.map((session) => (
                        <li
                            key={session.id}
                            className={
                                session.id === currentSessionId
                                    ? "conversation-item active"
                                    : "conversation-item"
                            }
                        >
                            {editingId === session.id ? (

                                <input
                                    autoFocus
                                    className="conversation-rename-input"
                                    value={editingTitle}
                                    onChange={(event) =>
                                        setEditingTitle(event.target.value)
                                    }
                                    onBlur={() => commitRename(session.id)}
                                    onKeyDown={(event) => {

                                        if (event.key === "Enter") {
                                            event.currentTarget.blur();
                                        }

                                        if (event.key === "Escape") {
                                            setEditingId(null);
                                        }
                                    }}
                                />

                            ) : (

                                <button
                                    className="conversation-title"
                                    onClick={() => handleSelect(session.id)}
                                >
                                    {session.title}
                                </button>

                            )}

                            <button
                                className="conversation-menu-toggle"
                                onMouseDown={(event) => event.stopPropagation()}
                                onClick={(event) => toggleMenu(event, session.id)}
                            >
                                ⋮
                            </button>
                        </li>
                    ))}

                </ul>

            </div>


            <footer className="sidebar-footer">

                <button>
                    ⬆ Upload
                </button>

            </footer>

            {menu && activeMenuSession && createPortal(

                <div
                    className="conversation-menu"
                    style={{ top: menu.top, left: menu.left }}
                    onMouseDown={(event) => event.stopPropagation()}
                >

                    <button
                        onClick={() =>
                            startRename(activeMenuSession.id, activeMenuSession.title)
                        }
                    >
                        Rename
                    </button>

                    <button
                        onClick={() => handleDelete(activeMenuSession.id)}
                    >
                        Delete
                    </button>

                </div>,

                document.body,
            )}

        </aside>
    );
}
