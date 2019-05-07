import * as React from "react";
import { Card, Button } from "tabler-react";

const LevelCardPreview = () => {
    // const { title, author, description } = props;

    return (
        <Card isColapsible>
        <Card.Header>
          <Card.Title>Level Title</Card.Title>
          <Card.Options>
            <Button RootComponent="a" color="info" size="sm">
              See level
            </Button>
          </Card.Options>
        </Card.Header>
        <Card.Body>
          Level descriptions
        </Card.Body>
      </Card>
    )
}

export default LevelCardPreview;