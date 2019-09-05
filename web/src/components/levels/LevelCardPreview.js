/* eslint-disable react/prop-types */
import * as React from "react";
import { Link } from "gatsby";
import { Card, Button } from "tabler-react";

const LevelBody = (props) => {
    const { author, description, locale, key } = props;

    return (
        <React.Fragment key={key}>
            <strong><small>Author: </small>{author}</strong><br />
            <strong><small>Locale: </small>{locale}</strong>
            <br />
            <br />
            <p>{description}</p>
            <Button.List>
                <Button RootComponent={Link} to="/" color="info" size="sm">
                    View level
                </Button>
                <Button RootComponent={Link} to="/" color="success" size="sm">
                    Validate level
                </Button>
            </Button.List>
        </React.Fragment>
    )
}

const LevelCardPreview = (props) => {
    const { levels } = props;

    return levels.map((level) =>
    <Card title={level.name} key={level.metadata.id}
        isCollapsible
        statusColor="orange"
        body={<LevelBody author={level.author} description={level.description} locale={level.locale} />}
    />)
}


export default LevelCardPreview;
