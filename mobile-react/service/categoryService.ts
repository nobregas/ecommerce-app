import api from "./apiClient";

export interface Category {
    id: number;
    name: string;
    imageUrl: string;
}


class CategoryService {
    async getAllCategories(): Promise<Category[]> {
        try {
            const response = await api.get("/category");
            console.log("get All categories", response.status)

            return response.data
        } catch (error) {
            throw this.handleError(error);
        }
    }

    private handleError(error: any): Error {
        const defaultMessage = "Failed to fetch category data";
        if (error.message) {
            return new Error(error.message || defaultMessage);
        }
        return new Error(defaultMessage);
    }
}

export default new CategoryService();