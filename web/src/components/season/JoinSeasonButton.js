import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";
import {useDispatch} from "react-redux";
import {Autocomplete} from "@material-ui/lab";
import {TextField} from "@material-ui/core";
import {createTeam as createTeamAction} from "../../actions/seasons";

const styles = {
  display: "flex",
  flexDirection: "column",
  marginBottom: "0.5rem",
  alignItems: "center",
}

const joinSeasonButton = ({seasons}) => {
  const dispatch = useDispatch();
  const [isFormOpen, setFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [season, setSeason] = useState("");
  const [error, setError] = useState(false);
  const createTeam = (seasonID, organizationID, name) =>
    dispatch(createTeamAction(seasonID, organizationID, name));

  useEffect(() => {
    if (!isEmpty(name)) {
      setError(false);
    }
  }, [name]);

  const handleNameChange = event => setName(event.target.value);

  const handleFormOpen = event => {
    event.preventDefault();
    setFormOpen(!isFormOpen);
  };

  const submitCreateTeam = async event => {
    event.preventDefault();
    if (!isEmpty(season) && !isEmpty(name)) {
      await createTeam(season, "", name);
      setFormOpen(false);
    }
  }

  return (
    <div style={styles}>
      <Button
        color="primary"
        onClick={handleFormOpen}
        icon={"plus"}
        size="sm"
        css={styles}
      >
        <FormattedMessage id="HomePage.joinSeason" />
      </Button>
      {isFormOpen && (
        <form onSubmit={submitCreateTeam}>
          <Form.FieldSet css={styles}>
            <Form.Group isRequired label="Organization Name" css={styles}>
              <Form.Input
                name="name"
                onChange={handleNameChange}
                invalid={error}
                cross={error}
                feedback={error && "Please, insert an organization name"}
              />
            </Form.Group>
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
            <Form.Group>
              <Button type="submit" color="primary" className="ml-auto">
                <FormattedMessage id="InviteOrganizationMemberButton.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </div>
  );
}

export default joinSeasonButton;
