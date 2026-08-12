import type { Document } from "./document";
import type { Note } from "./note";

export type KnowledgeSelection =
    | { type: "note"; item: Note }
    | { type: "document"; item: Document };
