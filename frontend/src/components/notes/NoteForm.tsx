import { useState, type FormEvent } from "react";

import type { Note } from "@/types/note";

interface NoteFormProps {
    onCreated: (note: Omit<Note, "id">) => void;
}

export default function NoteForm({
    onCreated,
}: NoteFormProps) {

    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>
    ) {
        event.preventDefault();

        const trimmedTitle = title.trim();
        const trimmedContent = content.trim();

        if (!trimmedTitle || !trimmedContent) {
            return;
        }

        onCreated({
            title: trimmedTitle,
            content: trimmedContent,
        });

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
                Save Note
            </button>

        </form>
    );
}