import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";
import {createTeam as createTeamAction} from "../../actions/seasons";
import {useDispatch} from "react-redux";
import {Autocomplete} from "@material-ui/lab";
import {TextField} from "@material-ui/core";

const styles = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  marginBottom: "0.5rem",
}

const CreateTeamButton = ({organizationID, seasons}) => {
  const dispatch = useDispatch();
  const [isFormOpen, setFormOpen] = useState(false);
  const [season, setSeason] = useState("");
  const createTeam = (seasonID, organizationID, name) =>
    dispatch(createTeamAction(seasonID, organizationID, name));

  const handleFormOpen = event => {
    event.preventDefault();
    setFormOpen(!isFormOpen);
  };

  const submitCreateTeam = async event => {
    event.preventDefault();
    if (!isEmpty(season)) {
      await createTeam(season, organizationID, "");
      setFormOpen(false);
    }
  }

  return (
    <div style={styles}>
      <Button
        color="primary"
        onClick={handleFormOpen}
        icon={"users"}
        size="sm"
        css={{  display: "flex",
          flexDirection: "column",
          alignItems: "center",
          marginBottom: "0.5rem",
          width: "131px"
      }}
      >
        {"Create new team"}
      </Button>
      {isFormOpen && (
        <form onSubmit={submitCreateTeam}>
          <Form.FieldSet css={styles}>
            <Form.Group isRequired label="Season" css={styles}>
              <Autocomplete
                freeSolo
                autoComplete
                autoHighlight
                onInputChange={(event, newInputValue) => {setSeason(newInputValue)}}
                options={seasons ? seasons.filter(season => !season.team).map(season => season.season.slug) : []}
                css={{width: "213px", backgroundColor: "white", height: "38px"}}
                renderInput={(params) => (
                  <TextField {...params}
                             size="small"
                             variant="outlined"

                  />
                )}
              />
            </Form.Group>
            <Form.Group css={styles}>
              <Button type="primary" color="primary" className="ml-auto">
                <FormattedMessage id="CreateTeamButton.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </div>
  );
}

export default CreateTeamButton;
