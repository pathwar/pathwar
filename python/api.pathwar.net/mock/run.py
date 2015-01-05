from eve import Eve
from settings import DOMAIN

from seeds import load_seeds


app = Eve()


if __name__ == '__main__':
    # Initialize data
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

        load_seeds(app)

    # Run
    app.run(debug=True, host='0.0.0.0')
