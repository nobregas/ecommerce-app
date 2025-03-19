todos:
    validate cpf in registering
    auth for all endpoints
    role middleware
    all CRUDS
    can create a product with images
    pagination
    get products by category
    Refactor the product quantity/stock to be an atomic and concurrent safe. The current implementation might lead to invalid stock values if multiple requests are made at the same time. Not only that, every time a product quantity is updated we need to query the products table (violates database normalization principle)
    create order-history table
    create user-address table
