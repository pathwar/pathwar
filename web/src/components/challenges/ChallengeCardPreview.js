/* eslint-disable react/prop-types */
import * as React from "react";
import { Link } from "gatsby";
import { Card, Button, Dimmer } from "tabler-react";
import styles from "../../styles/layout/loader.module.css";

const ChallengeBody = (props) => {
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
                    View challenge
                </Button>
                <Button RootComponent={Link} to="/" color="success" size="sm">
                    Validate challenge
                </Button>
            </Button.List>
        </React.Fragment>
    )
}

const ChallengeCardPreview = (props) => {
    const { challenges } = props;

    return !challenges ? <Dimmer className={styles.dimmer} active loader /> : challenges.map((challenge) =>
    <Card title={challenge.name} key={challenge.id}
        isCollapsible
        statusColor="orange"
        body={<ChallengeBody author={challenge.author} description={challenge.description} locale={challenge.locale} />}
    />)
}


export default ChallengeCardPreview;
