import { SparkIcon, UserIcon } from "@/components/common/icons";

export default function MessageAvatar({
    role,
}: {
    role: "user" | "assistant";
}) {
    return (
        <span
            className={`message-avatar message-avatar-${role}`}
            aria-hidden="true"
        >
            {role === "assistant"
                ? <SparkIcon size={16} />
                : <UserIcon size={16} />}
        </span>
    );
}
