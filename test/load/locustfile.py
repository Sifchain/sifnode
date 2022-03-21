from locust import HttpUser, User, task

class SifnodeUser(User):
    sequence_id = 0

    expectedResults = {}

    @task
    def export_pegged_token_to_evm(self):
        # Submit grpc to export ceth
        print("Sending pegged token to evm")

    @task
    def export_ibc_token_to_evm(self):
        print("Sending ibc token to evm")