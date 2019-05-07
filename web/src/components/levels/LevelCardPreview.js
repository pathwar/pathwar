import * as React from "react";
import { Link } from "react-router-dom";
import { Card, Button } from "tabler-react";

const LevelBody = () => {
    return (
        <React.Fragment>
            <strong>Author: Author Name</strong>
            <p>Level description</p>
            <Button RootComponent={Link} to="/" color="info" size="sm">
                See level
            </Button>
        </React.Fragment>
    )
}

const LevelCardPreview = () => {
    // const { title, author, description } = props;

    return (
        <Card title="Level Title"
        isCollapsible
        statusColor="orange" 
        body={<LevelBody />}
        />
    )
}

export default LevelCardPreview;