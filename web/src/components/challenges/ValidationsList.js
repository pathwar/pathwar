import React from "react"
import { Card, Grid } from "tabler-react"

const ValidationsList = () => {
  return (
    <Grid.Row cards={true}>
      <Grid.Col lg={4} md={4} sm={6} xs={6}>
        <Card title="Validation" statusColor="orange" isCollapsible>
          <Card.Body>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aperiam
            deleniti fugit incidunt, iste, itaque minima neque pariatur
            perferendis sed suscipit velit vitae voluptatem. A consequuntur,
            deserunt eaque error nulla temporibus!
          </Card.Body>
        </Card>
      </Grid.Col>

      <Grid.Col lg={4} md={4} sm={6} xs={6}>
        <Card title="Validation" statusColor="orange" isCollapsible>
          <Card.Body>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aperiam
            deleniti fugit incidunt, iste, itaque minima neque pariatur
            perferendis sed suscipit velit vitae voluptatem. A consequuntur,
            deserunt eaque error nulla temporibus!
          </Card.Body>
        </Card>
      </Grid.Col>
    </Grid.Row>
  )
}

export default ValidationsList
