import { useEffect, useState, type FormEvent } from "react";

import type { Note } from "@/types/note";

interface NoteFormProps {
    selectedNote?: Note | null;
    onCreated: (note: Omit<Note, "id">) => void;
    onUpdated: (id: number, note: Omit<Note, "id">) => void;
}

export default function NoteForm({
    selectedNote,
    onCreated,
    onUpdated,
}: NoteFormProps) {

    const [title, setTitle] = useState(selectedNote?.title ?? "");
    const [content, setContent] = useState(selectedNote?.content ?? "");

    useEffect(() => {

        setTitle(selectedNote?.title ?? "");

        setContent(selectedNote?.content ?? "");

    }, [selectedNote]);

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>
    ) {
        event.preventDefault();

        const trimmedTitle = title.trim();
        const trimmedContent = content.trim();

        if (!trimmedTitle || !trimmedContent) {
            return;
        }

        if (selectedNote) {

            onUpdated(
                selectedNote.id,
                {
                    title: trimmedTitle,
                    content: trimmedContent,
                }
            );

        } else {

            onCreated({
                title: trimmedTitle,
                content: trimmedContent,
            });

        }

        setTitle("");
        setContent("");
    }

    return (
        <form onSubmit={handleSubmit}>

            <input
                type="text"
                placeholder="Title"
                value={title}
                onChange={(event) =>
                    setTitle(event.target.value)
                }
            />

            <textarea
                placeholder="Content"
                value={content}
                onChange={(event) =>
                    setContent(event.target.value)
                }
            />

            <button type="submit">

                {
                    selectedNote
                        ? "Update Note"
                        : "Save Note"
                }

            </button>

        </form>
    );
}
