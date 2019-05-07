/* eslint-disable react/prop-types */
import * as React from "react";
import { Link } from "react-router-dom";
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
            <Button RootComponent={Link} to="/" color="info" size="sm">
                See level
            </Button>
        </React.Fragment>
    )
}

const LevelCardPreview = (props) => {
    const { levels } = props;

    return levels.map((level) => 
    <Card title={level.name} key={level.metadata.id}
        isCollapsible
        statusColor="orange" 
        body={<LevelBody key={level.metadata.id} author={level.author} description={level.description} locale={level.locale} />}
    />)
}


export default LevelCardPreview;