import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import EmptyState from "@/components/common/EmptyState";

describe("EmptyState", () => {

    it("renders the icon, title, description and actions", () => {
        render(
            <EmptyState
                icon={<span data-testid="icon" />}
                title="Nothing here"
                description="Add something first."
                actions={<button type="button">Add</button>}
            />,
        );

        expect(screen.getByTestId("icon")).toBeInTheDocument();
        expect(screen.getByText("Nothing here")).toBeInTheDocument();
        expect(screen.getByText("Add something first.")).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Add" })).toBeInTheDocument();
    });

    it("omits the description and actions when they are not given", () => {
        render(<EmptyState icon={<span />} title="Just a title" />);

        expect(screen.getByText("Just a title")).toBeInTheDocument();
        expect(screen.queryByRole("button")).not.toBeInTheDocument();
    });
});
